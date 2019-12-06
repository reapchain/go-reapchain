#!/bin/bash
#!/usr/bin/python

read -p "Execute script:(y/n) " response
if [ "$response" = "y" ]; then
    echo -e "\n\nLoading....\n\n"

    for ((x = 0; x<5; x++))
    do
        echo -e "Open $x terminal\n\n"
        open -na Terminal.app
    done
fi