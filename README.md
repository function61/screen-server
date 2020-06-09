![Build status](https://github.com/function61/screen-server/workflows/Build/badge.svg)
[![Download](https://img.shields.io/docker/pulls/fn61/screen-server.svg?style=for-the-badge)](https://hub.docker.com/r/fn61/screen-server/)

Minimal VNC-servable desktop environment with a web browser running in a Alpine Linux
Docker container, for displaying a webpage on an untrusted device like an old Android tablet.

Nice added benefit is that you can connect to the same screen simultaneously from other
devices like PCs as well, and if scripting is needed (show content X on screen Y), it's
easier to achieve things like this on a PC than on a tablet.

![](docs/network-drawing.png)


Why?
----

I had old Android tablets lying around. They are untrusted, because they are dangerous,
because they haven't received software updates. I wanted to use one of them as always-on
info screen.

Because they're untrusted, I:

- taped off the camera & mic
- put it in a guest Wifi network (with no access to my LAN)
- configured firewall to not let the tablet in the internet, but only access to VNC port
  in my LAN

(Pro-tip: one of my tablets was so old (Android 4.2) that it didn't even support
[RealVNC's current Android package](https://play.google.com/store/apps/details?id=com.realvnc.viewer.android).
I downloaded an old version via apkpure.com, and firewalled the tablet off before beginning
the installation.)


What's special about this image
-------------------------------

- Supports multiple screens (ran as separate users) with different resolutions
- OSD notifications, OSD API for sending messages to screens over the network
- Web UI w/ screen previews
- Small size (for an image with desktop environment and a web browser)


How to run
----------

```console
$ docker run -d \
	--name screen-server \
	-p 5900:5900 \
	-p 80:80 \
	-e "SCREEN_1=5900,800,1280,Galaxy Tab 2" \
	fn61/screen-server:TAG
```

The format for the `SCREEN_1` parameter is VNC port, display width, height, screen name
(web UI shows this, some VNC clients show this)

If you have more than one screen, just add `SCREEN_2` and so on..


Note about state
----------------

Users should assume all state gets wiped daily. It doesn't, but we have absolutely no plans
to support migrating state (like Firefox configuration or installed plugins) when new
versions of this image gets released. When you spin up a new container with the new version,
all state gets lost.


Web UI
------

It shows you the preview on what's all the screens.

![](docs/web-ui.png)

TODO: maybe add [noVNC](https://github.com/novnc/noVNC).


Sending OSD notifications
-------------------------

My use case was to display messages sent by my home automation in the screen that is always
visible.

```console
$ curl -d "msg=Hello world" http://localhost/api/screen/1/osd/notify
```

The notification is visible for a few seconds.

We have plugin drivers for different OSD notification implementations -  currently we use
[zenity](https://en.wikipedia.org/wiki/Zenity), which is not pretty. A prettier way could
be to show the notification as a full-screen webpage (so we get CSS animations etc.), but
that's still TODO.

![](docs/osd-notification.png)


Credits
-------

I learned how to plug Xvfb, x11vnc, Openbox together from
[danielguerra69/alpine-vnc](https://github.com/danielguerra69/alpine-vnc). I added Firefox,
process management in Go (instead of Python), multi-screen support, OSD etc.
