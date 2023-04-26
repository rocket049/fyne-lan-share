fyneLanShare:*.go
	go build -ldflags -s
install:fyneLanShare fyneLanShare.svg fyneLanShare.desktop
	install -m 755 -s -t /usr/local/bin fyneLanShare
	install -t /usr/share/icons fyneLanShare.svg
	install -t /usr/share/applications fyneLanShare.desktop
uninstall:
	rm -f /usr/local/bin/fyneLanShare /usr/share/icons/fyneLanShare.svg /usr/share/applications/fyneLanShare.desktop
