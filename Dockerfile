FROM ubuntu:latest

RUN apt update && apt install -y xvfb openbox xfce4-terminal x11vnc dbus ca-certificates firefox

ADD misc/menu.xml /etc/xdg/openbox/

CMD ["/usr/local/bin/screen-server", "run"]

ADD rel/screen-server_linux-amd64 /usr/local/bin/screen-server
