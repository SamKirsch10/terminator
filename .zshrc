# If you come from bash you might have to change your $PATH.
export PATH=$HOME/go/bin:$HOME/bin:/usr/local/bin:/usr/local/go/bin:$PATH

# Apache Maven Environment Variables
# MAVEN_HOME for Maven 1 - M2_HOME for Maven 2
export M2_HOME=/usr/local/src/apache-maven
export PATH=${M2_HOME}/bin:${PATH}

# Path to your oh-my-zsh installation.
  export ZSH="/home/skirsch/.oh-my-zsh"

# Set name of the theme to load. Optionally, if you set this to "random"
# it'll load a random theme each time that oh-my-zsh is loaded.
# See https://github.com/robbyrussell/oh-my-zsh/wiki/Themes
ZSH_THEME="bullet-train"


# Uncomment the following line to disable bi-weekly auto-update checks.
DISABLE_AUTO_UPDATE="true"

# Which plugins would you like to load? (plugins can be found in ~/.oh-my-zsh/plugins/*)
# Custom plugins may be added to ~/.oh-my-zsh/custom/plugins/
# Example format: plugins=(rails git textmate ruby lighthouse)
# Add wisely, as too many plugins slow down shell startup.
plugins=(
  git
  colored-man-pages
  colorize
  autojump
  docker
  gitfast
)

source $ZSH/oh-my-zsh.sh

export ANSIBLE_VAULT_PASSWORD_FILE=/home/skirsch/v

# alias venv="source ~/venv/bin/activate"
alias venv="source ~/venv/bin/activate"
alias venv3="source ~/venv3/bin/activate"

alias master="git checkout master && git pull"
alias rm_windows_line_ends="sed -i -e 's/\r$//'"
alias tempdir="cd $(mktemp -d)"
alias passhash="openssl passwd -1"
alias socksProxy="/home/skirsch/repos/samkirsch10-terminator/scripts/socksProxy.sh"
alias dfimage="docker run -v /var/run/docker.sock:/var/run/docker.sock --rm dduvnjak/dockerfile-from-image"
eval $(thefuck --alias)
alias retry-ssh="~/repos/samkirsch10-terminator/scripts/retry-ssh.sh"

## KUBERNETES STUFFs
alias getNodePorts="kubectl get svc --all-namespaces -o go-template='{{range .items}}{{range.spec.ports}}{{if .nodePort}}{{.nodePort}} - {{.name}}{{\"\n\"}}{{end}}{{end}}{{end}}' "
alias show-contexts="kubectl config get-contexts"
alias k8s-switch=/home/skirsch/repos/samkirsch10-terminator/scripts/kubeSwitch.sh

##AWS Stuff
alias aws-profile="/home/skirsch/repos/samkirsch10-terminator/scripts/awsSetProfile.sh"



eval `ssh-agent` > /dev/null 2>&1
ssh-add ~/.ssh/skirsch_1537968652  > /dev/null 2>&1

nohup ~/killstuff.sh > /dev/null 2>&1 &
