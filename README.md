# rightdog

Earn money like a dog, spend it like a gentleman!! Gazaaaaaaaaaa

Docker Installation
===================

        curl -sSL https://get.docker.com | sh

Golang Installation
===================

        curl -L https://raw.github.com/grobins2/gobrew/master/tools/install.sh | sh
        export PATH="$HOME/.gobrew/bin:$PATH"
        eval "$(gobrew init -)"
        mkdir ~/go && cd ~/go
        gobrew workspace set
        gobrew install 1.9.2
        gobrew use 1.9.2
