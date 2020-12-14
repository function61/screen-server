FROM alpine:latest

# - dbus and ttf-freefont dependencies of Firefox
# - zenity for OSD notifications
# - scrot for screenshots
RUN addgroup alpine && apk add --update \
	xvfb \
	openbox \
	xfce4-terminal \
	x11vnc \
	dbus \
	ttf-freefont \
	ca-certificates \
	firefox-esr \
	scrot \
	zenity

ADD misc/menu.xml /etc/xdg/openbox/

CMD ["/usr/local/bin/screen-server", "run"]

ADD rel/screen-server_linux-amd64 /usr/local/bin/screen-server
