# GoToShock

Reimplementation of a well-known device to allow shocking your friends and/or loved ones via the internet - or local network or whatever.

Basically:

```
while true; do curl -v 'http://raspberrypi:8080/v1alpha1/message/hellorld!/1/vibrate/100' -X POST; sleep 2; done
```

:>
