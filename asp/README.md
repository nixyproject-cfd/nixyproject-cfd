## UPGRADE FOR DEBIAN
```
apt update -y && apt upgrade -y && apt dist-upgrade -y && reboot
```

##  UPGRADE FOR UBUNTU
```
apt update && apt upgrade -y && update-grub && sleep 2 && reboot
```

## INSTALL SCRIPT 
```
apt install -y && apt update -y && apt upgrade -y && wget -q https://raw.githubusercontent.com/nixyproject-cfd/asp/main/ubu20-deb10-stable.sh && chmod +x ubu20-deb10-stable.sh && ./ubu20-deb10-stable.sh
```

## UPDATE SCRIPT
```
wget -q https://raw.githubusercontent.com/nixyproject-cfd/asp/main/update.sh && chmod +x update.sh && ./update.sh
```
