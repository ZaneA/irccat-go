irccat.go
===

My first ever golang program, woo! A small utility for posting things to IRC.

Building
---

`go build irccat.go`

or to force 32bit (so I can blast the result to my VPS)

`GOARCH=386 go build irccat.go`

Running
---

    $ irccat -h
    Usage of irccat:
      -dest="": Where to send the text (nickname or channel)
      -nick="irccat": Nickname
      -server="": Server to send to (e.g. chat.freenode.net:6667)
      -verbose=false: Be verbose (and stay in the foreground)
    $ echo Hello World | irccat -dest='#irccat-test' -server=chat.freenode.net:6667 -verbose
     ... a bunch of output here and hopefully a message in your channel ...

License
---

Not that it's large enough to warrant one... I hereby place this code in the public domain.
