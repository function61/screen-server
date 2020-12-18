FROM ubuntu:latest

RUN apt update && DEBIAN_FRONTEND=noninteractive apt install -y \
	xvfb \
	x11vnc \
	openbox \
	xfce4-terminal \
	dbus \
	ca-certificates \
	firefox

# install noVNC
RUN mkdir -p www/vnc \
	&& cd www/vnc \
	&& curl -L https://github.com/novnc/noVNC/archive/v1.2.0.tar.gz | tar --strip-components=1 -xz noVNC-1.2.0/app noVNC-1.2.0/core noVNC-1.2.0/vendor noVNC-1.2.0/vnc.html \
	&& mv vnc.html index.html

ADD misc/menu.xml /etc/xdg/openbox/

CMD ["/usr/local/bin/screen-server", "run"]

ADD rel/screen-server_linux-amd64 /usr/local/bin/screen-server
