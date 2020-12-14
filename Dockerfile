FROM ubuntu:latest

RUN apt update && DEBIAN_FRONTEND=noninteractive apt install -y \
	xvfb \
	x11vnc \
	openbox \
	xfce4-terminal \
	dbus \
	ca-certificates \
	firefox

ADD misc/menu.xml /etc/xdg/openbox/

CMD ["/usr/local/bin/screen-server", "run"]

ADD rel/screen-server_linux-amd64 /usr/local/bin/screen-server
