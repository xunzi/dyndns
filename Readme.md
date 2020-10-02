# Dyndns - self-hosted dynamic dns pointer to an arbitrary dns name

This is meant to be run from cron every hour or so (the interval depends on
how often your ip address changes).

The program tries to fetch its public ip from an API endpoint (this will be made configurable) and
compares the ip to the dns A record you supply on the command line. If the ips differ the program
will update the DNS record (currently only works with Hetzer's DNS API, will also be configurable)
to point to the new address.

It is probably not very useful at the moment but only my attempt to build something usable while learning
Golang.
