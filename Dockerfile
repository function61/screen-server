FROM ubuntu:latest

# can't have default values here, otherwise they'd overwrite the buildx-supplied ones
ARG TARGETOS
ARG TARGETARCH

RUN apt update && DEBIAN_FRONTEND=noninteractive apt install -y \
	xvfb \
	x11vnc \
	openbox \
	xfce4-terminal \
	dbus \
	curl \
	ca-certificates \
	software-properties-common

# to add non-snap Firefox: https://askubuntu.com/questions/1399383/how-to-install-firefox-as-a-traditional-deb-package-without-snap-in-ubuntu-22
#
# the Canonical folks make miserable decisions and tell people when they try to use apt to install Firefox,
# they actually tell "no you are not allowed to do that" and just install a stub to tell people to use
# their Snap instead. so that results in Ubuntu breaking Firefox for containerized installs..
# so I'm supposed to run Snap in a Docker container? imagine being so incompetent to make decisions like this at Canonical.
#
# Firefox in Snap doesn't even work properly. people's printers etc. settings seem to have broken.
# imagine being ok with trying to force people to use an inferior variant of package.
RUN add-apt-repository ppa:mozillateam/ppa

RUN echo "Package: *\nPin: release o=LP-PPA-mozillateam\nPin-Priority: 1001\n\nPackage: firefox\nPin: version 1:1snap1-0ubuntu2\nPin-Priority: -1\n" > /etc/apt/preferences.d/mozilla-firefox

RUN cat /etc/apt/preferences.d/mozilla-firefox && apt update && DEBIAN_FRONTEND=noninteractive apt install -y \
	firefox

# install noVNC
RUN mkdir -p www/vnc \
	&& cd www/vnc \
	&& curl -L https://github.com/novnc/noVNC/archive/v1.2.0.tar.gz | tar --strip-components=1 -xz noVNC-1.2.0/app noVNC-1.2.0/core noVNC-1.2.0/vendor noVNC-1.2.0/vnc.html \
	&& mv vnc.html index.html

ADD misc/menu.xml /etc/xdg/openbox/

# the openbox config file (rc.xml) is a file that pretty much has to list all behaviour of openbox,
# so instead of having our one config there it is around 1 000 lines so the least bad way for us to
# make Firefox maximized is to add the application rule via search-and-replace (and not e.g. copy snapshot
# of the config file and add our change in there) in order to get future updates to that important file.
RUN sed -i '/<applications>/a\
    <application class="firefox">\
        <maximized>true</maximized>\
    </application>' /etc/xdg/openbox/rc.xml

CMD ["/usr/local/bin/screen-server", "run"]

ADD rel/screen-server_linux-$TARGETARCH /usr/local/bin/screen-server

USER ubuntu
