# Usage

For use with nagios, sensu or some other monitoring software that supports
external programs.

	$ resque-ok -h
	Usage: resque-ok [options] QUEUE [...]
	  -ns="resque": reqsue namespace

Call passing one or more resque queue names as arguments. It also takes an
optional namespace argument that should coincide with your
Resque.redis.namespace setting.

## Comment

It compares the top of the queue against a saved value and if they are the same
then the queue isn't getting processes. If the queue is empty or the value from
the queue is different than the saved value, all is ok.

Note that redis doesn't keep entries for empty queues, so there is no way to
tell if the passed in strings are the correct queue names. So double check.

Obviously make sure the queue is processes more frequently than you run this.

## License

MIT licence. See LICENCE file.

