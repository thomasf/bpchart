
# charts and web service for omron bp monitors

*Very much a work in progress atm, just putting it in gh so I don't loose my WIP*

Some screenshots using generated test data:

Averaged view:

![one](https://raw.github.com/thomasf/bpchart/master/bpchart01.png)

All entries view:

![two](https://raw.github.com/thomasf/bpchart/master/bpchart02.png)

## pre requirements

requires `libomron` to be installed. (will be migrated to using a libusb go wrapper instead)

## Setting up usb permissions under ubuntu/udev

Create udev rule in `/etc/udev/rules.d/51-omron.rules` with content:

```
SUBSYSTEM=="usb", ATTR{idVendor}=="0590", MODE="0664", GROUP="plugdev"
```

If it's not working right away, force udev to be reloaded:

```sh
sudo udevadm control --reload
sudo udevadm trigger
```

You also need to be in the `plugdev` group .
