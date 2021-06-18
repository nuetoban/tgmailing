# tgmailing

✉️ **tgmailing** is a tool which can help you with sending notifications to users of your Telegram Bot.

## Requirements

To run this program you should have:
1. Ad post.
2. List of bots.
3. List of chats which will recieve the post.
4. Service chat which contain all specified bots.
   All the bots should be presented in this chat.
   It's needed if you send post with media.
5. A chat which will recieve service notifications (mailing start/finish, etc.). Not required
   unless you use `--start-notification`, `--finish-notification`, `--each-finish-notification` flags.

## Installation
```bash
go install github.com/nuetoban/tgmailing
```

## Usage

The following example will take chats from `input_examples/chats.lines` file,
ad post from `input_examples/example.json` and bots from `input_examples/bots.lines`.
```bash
tgmailing \
  --ad-file           input_examples/example.json \
  --chats-file        input_examples/chats.lines \
  --bots-file         input_examples/bots.lines \
  --service-chat      -598317757 \
  --notification-chat -598317757 \
  --start-notification \
  --finish-notification
```

### PostgreSQL
If you have PostgreSQL, you can use `PGSQL` input source.
```bash
tgmailing \
  --chats-src         PGSQL \
  --chats-query       'SELECT id FROM chats' \
  --ad-file           input_examples/example.json \
  --bots-file         input_examples/bots.lines \
  --service-chat      -598317757 \
  --notification-chat -598317757 \
  --start-notification \
  --finish-notification
```

### Advice: fill chats list

If you want to get all chats for your bot and you use, for example, PostgreSQL with table `chats`,
which contain a column `id`, you can fill `chats.lines` file via following psql command:
```bash
psql -U postgres -t -c "SELECT id FROM chats" > input_examples/chats.lines
```

### Help
```
usage: tgmailing [-h|--help] [--ad-src (JSONFILE|PGSQL)] [--bots-src
                 (ENV|LINESFILE|PGSQL)] [--chats-src (LINESFILE|PGSQL)]
                 [--ad-file "<value>"] [--bots-file "<value>"] [--chats-file
                 "<value>"] [--bots-query "<value>"] [--chats-query "<value>"]
                 [--ad-query "<value>"] [--bots-db-env-prefix "<value>"]
                 [--chats-db-env-prefix "<value>"] [--ad-db-env-prefix
                 "<value>"] [--no-server] [-m|--metrics-port <integer>]
                 [-n|--start-notification] [-f|--finish-notification]
                 [--each-finish-notification] [--notification-chat <integer>]
                 -s|--service-chat <integer>

                 Sends post to Telegram chats via Bots

Arguments:

  -h  --help                      Print help information
      --ad-src                    Ad source. Default: JSONFILE
      --bots-src                  Bots source. Default: LINESFILE
      --chats-src                 Chats source. Default: LINESFILE
      --ad-file                   Path to Ad file
      --bots-file                 Path to bots file
      --chats-file                Path to chats file
      --bots-query                SQL query to fetch Bots
      --chats-query               SQL query to fetch Chats
      --ad-query                  SQL query to fetch Ad post
      --bots-db-env-prefix        Prefix for DB env credentials for Bots.
                                  Default: SENDER_
      --chats-db-env-prefix       Prefix for DB env credentials for Chats.
                                  Default: SENDER_
      --ad-db-env-prefix          Prefix for DB env credentials for Ad post.
                                  Default: SENDER_
      --no-server                 Disable metrics server
  -m  --metrics-port              Metrics server port. Default: 9090
  -n  --start-notification        Send message to chat on start
  -f  --finish-notification       Send message to chat on finish
      --each-finish-notification  Send message to chat on finish for each bot
      --notification-chat         Chat to send notifications
  -s  --service-chat              Chat to send files
```

