# BENCHmarker for Icinga DB's Redis

```
$ redis-cli flushdb
OK
$ go run . -host localhost -port 6380
0.00  wip: [*]
0.00  done: [command icingaapplication logger perfdatawriter filelogger apilistener sysloglogger dependency comment scheduleddowntime gelfwriter icingadb apiuser downtime notification service idomysqlconnection graphitewriter configobject application timeperiod elasticsearchwriter customvarobject servicegroup idopgsqlconnection user endpoint streamlogger influxdbwriter dbconnection usergroup eventcommand zone]
0.04  done: [notificationcommand notificationcomponent opentsdbwriter]
0.54  done: [checkcommand checkable checkercomponent]
57.67  done: [host]
57.71  done: [hostgroup *]
67.89
```

`0.00  wip: [*]` means Icinga 2 has started writing.

`57.67  done: [host]` means Icinga 2 has finished writing hosts
57.67s after it has started writing.

`57.71  done: [*]` means Icinga 2 has finished writing everything
57.71s after it has started writing.
