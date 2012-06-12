# Usage

Intended to be used to monitor Redis [http://redis.io/] queues. Originally
developed to monitor resque [https://github.com/defunkt/resque], it now
supports any redis queue.  It can run in standalone mode (cron/monit) where it
sends emails on alert or it can work as a plugin for nagios, sensu and possibly
other monitoring software that supports external programs.

    $ queue-ok -h
    Usage: [options] ./queue-ok
    Options (for optional email message):
      -e=false: enable email message
      -f="": From: address
      -l=false: left push (test index -1)
      -p="resque*:queue:*": queue key pattern
      -r=true: right push (test index 0)
      -s="localhost:25": smtp server
      -t="": To: address
    Returns:
    0 on success
    1 not used
    2 when queue issued not being processed
    3 when there is an error with the check

## How it works

It looks for any redis keys that follow the pattern of a queue and monitors
their first/last entry (rpush/lpop vs lpush/rpop). If it doesn't change between
2 calls an alert is issued (email sent/return code 2).

If the queue is empty or the value from the queue is different than the saved
value, all is ok. Make sure the queue is processes more frequently than you run
this.

## License

MIT licence. See LICENCE file.

