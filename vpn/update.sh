#!/bin/bash
clear
fun_bar() {
    CMD[0]="$1"
    CMD[1]="$2"
    (
        [[ -e $HOME/fim ]] && rm $HOME/fim
        ${CMD[0]} -y >/dev/null 2>&1
        ${CMD[1]} -y >/dev/null 2>&1
        touch $HOME/fim
    ) >/dev/null 2>&1 &
    tput civis
    echo -ne "  \033[0;33mPlease Wait Loading \033[1;37m- \033[0;33m["
    while true; do
        for ((i = 0; i < 18; i++)); do
            echo -ne "\033[0;32m#"
            sleep 0.1s
        done
        [[ -e $HOME/fim ]] && rm $HOME/fim && break
        echo -e "\033[0;33m]"
        sleep 1s
        tput cuu1
        tput dl1
        echo -ne "  \033[0;33mPlease Wait Loading \033[1;37m- \033[0;33m["
    done
    echo -e "\033[0;33m]\033[1;37m -\033[1;32m OK !\033[1;37m"
    tput cnorm
}
res1() {
    wget https://raw.githubusercontent.com/nixyproject-cfd/asp/main/menu/menu.zip
    wget -q -O /usr/bin/enc "https://raw.githubusercontent.com/nixyproject-cfd/asp/main/encrypt/epro" ; chmod +x /usr/bin/enc
    7z e -paskykenza123 menu.zip
    unzip menu.zip
    chmod +x menu/*
    enc menu/*
    mv menu/* /usr/local/sbin    
    #unzip menu.zip
    #chmod +x menu/*
    #mv menu/* /usr/local/sbin
    rm -rf menu.zip
    rm -rf update.sh
    wget -qO- fv-tunnel "https://raw.githubusercontent.com/nixyproject-cfd/asp/main/config/fv-tunnel" 
    chmod +x fv-tunnel 
    bash fv-tunnel
    rm -rf fv-tunnel
}
netfilter-persistent
clear
echo -e "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" | lolcat
echo -e " \e[1;97;101m            UPDATE SCRIPT               \e[0m"
echo -e "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" | lolcat
echo -e ""
echo -e "  \033[1;91m Update Script Service\033[1;37m"
fun_bar 'res1'
echo -e "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" | lolcat
echo -e ""
read -n 1 -s -r -p "Press [ Enter ] to back on menu"
menu
