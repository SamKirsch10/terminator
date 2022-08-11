# If you come from bash you might have to change your $PATH.
export PATH=$HOME/.local/bin/:$HOME/go/bin:$HOME/bin:/usr/local/bin:/usr/local/go/bin:$PATH

# Apache Maven Environment Variables
# MAVEN_HOME for Maven 1 - M2_HOME for Maven 2
export M2_HOME=/usr/local/src/apache-maven
export PATH=${M2_HOME}/bin:${PATH}

# Path to your oh-my-zsh installation.
export ZSH="$HOME/.oh-my-zsh"

export GOROOT=/usr/local/go

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
  docker
  gitfast
)

export ZSH_DISABLE_COMPFIX="true"
source $ZSH/oh-my-zsh.sh

export ANSIBLE_VAULT_PASSWORD_FILE=/home/skirsch/v

bash ~/repos/samkirsch10-terminator/scripts/keyExpire.sh

export DOCKER_BUILDKIT=0
source ~/venv3/bin/activate

# Random Alias stuff
alias python=python3
alias venv3="source ~/venv3/bin/activate"

alias rm_windows_line_ends="sed -i -e 's/\r$//'"
alias tempdir="cd $(mktemp -d)"
alias passhash="openssl passwd -1"
alias socksProxy="~/repos/github.com/samkirsch10-terminator/scripts/socksProxy.sh"
alias dfimage="docker run -v /var/run/docker.sock:/var/run/docker.sock --rm dduvnjak/dockerfile-from-image"
alias retry-ssh="~/repos/github.com/samkirsch10-terminator/scripts/retry-ssh.sh"
if command -v fuck &> /dev/null; then
    eval $(thefuck --alias)
fi
alias pssh="$HOME/go/bin/orgalorg -y"
alias ogrep="~/repos/github.com/samkirsch10/samkirsch10-terminator/scripts/ogrep.sh"
alias pfx-to-pem="~/repos/github.com/samkirsch10/samkirsch10-terminator/scripts/pfx-to-pem.sh"
alias subl="/Applications/Sublime\ Text.app/Contents/SharedSupport/bin/subl"
dexec() {
  id=$(docker ps | tail -1 | awk '{print $1}')
  docker exec -it $id $1
}

# git stuff
alias master="git checkout master && git pull"
alias remove-branches="git branch --merged | grep -v master | xargs git branch -d"
master() {
  default_branch=$(git remote show origin | grep 'HEAD branch' | cut -d' ' -f5)
  echo "git checkout $default_branch && git pull"
  git checkout $default_branch && git pull
}

# KUBERNETES STUFFs
alias getNodePorts="kubectl get svc --all-namespaces -o go-template='{{range .items}}{{range.spec.ports}}{{if .nodePort}}{{.nodePort}} - {{.name}}{{\"\n\"}}{{end}}{{end}}{{end}}' "
alias show-contexts="kubectl config get-contexts"
alias test-pod="kubectl run sxk3161-test-shell --rm -i --tty --image ubuntu -- bash"


# GCLOUD Stuff
# /Users/sxk3161/repos/github.com/samkirsch10/samkirsch10-terminator/scripts/iTermProdWatcher.py &
export USE_GKE_GCLOUD_AUTH_PLUGIN=True
g-proj() {
  echo "$(cat ~/.config/gcloud/configurations/config_default | grep project | awk '{print $NF}')"
}

alias gcp-switch="~/repos/github.com/samkirsch10/samkirsch10-terminator/scripts/gprojects.sh"
alias cls="gcloud container clusters list"
alias force_drain="~/repos/github.com/samkirsch10/samkirsch10-terminator/scripts/k8s_drain.sh"


