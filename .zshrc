# If you come from bash you might have to change your $PATH.
# export PATH=$HOME/bin:/usr/local/bin:$PATH

# Path to your oh-my-zsh installation.
export ZSH=$HOME/.oh-my-zsh

# Set name of the theme to load --- if set to "random", it will
# load a random theme each time oh-my-zsh is loaded, in which case,
# to know which specific one was loaded, run: echo $RANDOM_THEME
# See https://github.com/ohmyzsh/ohmyzsh/wiki/Themes
ZSH_THEME="bullet-train"


source $ZSH/oh-my-zsh.sh

plugins=(git colored-man-pages colorize pip python brew osx zsh-syntax-highlighting zsh-autosuggestions)


# Aliases
alias venv="source ~/venv/bin/activate"
alias python="python3"
alias subl="/Applications/Sublime\ Text.app/Contents/SharedSupport/bin/subl"
alias k8s-switch="~/repos/terminator/scripts/kubeSwitch.sh"

# eval `ssh-agent` > /dev/null 2>&1
# ssh-add ~/.ssh/id_rsa  > /dev/null 2>&1


nohup ~/repos/terminator/scripts/killstuff.sh > /dev/null 2>&1 &
