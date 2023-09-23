# If you come from bash you might have to change your $PATH.
export PATH=/opt/homebrew/bin:$HOME/go/bin:$HOME/bin:/usr/local/bin:/usr/local/go/bin:$PATH

# Apache Maven Environment Variables
# MAVEN_HOME for Maven 1 - M2_HOME for Maven 2
export M2_HOME=/usr/local/src/apache-maven
export PATH=${M2_HOME}/bin:${PATH}

# Path to your oh-my-zsh installation.
export ZSH="$HOME/.oh-my-zsh"

eval $(brew shellenv)

export PATH=$HOME/src/fsdev/tools/node/bin/:$HOME/src/fsdev/tools/python3/bin:$PATH
export PATH=$HOME/.local/bin:$PATH

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

# export ANSIBLE_VAULT_PASSWORD_FILE=/home/skirsch/v

export DOCKER_BUILDKIT=0
# source ~/venv3/bin/activate

# Random Alias stuff
alias python=python3
alias venv3="source ~/venv3/bin/activate"
monitor-input() {
  # https://github.com/waydabber/m1ddc
  BINARY="/Users/samkirsch/src/m1ddc/m1ddc"
  CMD="set input"
  INPUT="17" #default 17 = HDMI
  
  arg="echo $1 | tr '[:upper:]' '[:lower:]'" #toLower
  if [[ "$arg" == *"-h"* ]]; then
    echo "Usage:"
    echo "\tmonitor-input [input]"
    echo "\t[input] can be 'hdmi, dp, usbc'"
    return
  fi
  case "$arg" in
    hdmi)
      INPUT="17"
      ;;
    dp|displayport)
      INPUT="15"
      ;;
    usbc|usb-c|usb)
      INPUT="27"
      ;;
  esac
  echo "${BINARY} ${CMD} ${INPUT}"
  eval "${BINARY} ${CMD} ${INPUT}"
}

alias rm_windows_line_ends="sed -i -e 's/\r$//'"
alias tempdir="cd $(mktemp -d)"
alias passhash="openssl passwd -1"
alias socksProxy="~/src/terminator/scripts/socksProxy.sh"
alias dfimage="docker run -v /var/run/docker.sock:/var/run/docker.sock --rm dduvnjak/dockerfile-from-image"
alias retry-ssh="~/src/terminator/scripts/retry-ssh.sh"
if command -v fuck &> /dev/null; then
    eval $(thefuck --alias)
fi
alias pssh="$HOME/go/bin/orgalorg -y"
alias ogrep="~/src/terminator/scripts/ogrep.sh"
alias pfx-to-pem="~/src/terminator/scripts/pfx-to-pem.sh"
alias subl="/Applications/Sublime\ Text.app/Contents/SharedSupport/bin/subl"
alias re-search="~/src/terminator/scripts/regex-search.py"
dexec() {
  id=$(docker ps | tail -1 | awk '{print $1}')
  docker exec -it $id $1
}

# git stuff
# alias master="git checkout master && git pull"
remove-branches() {
  REMOTES=$(git ls-remote $(git remote -v | tail -1 | awk '{print $2}'))
  if [[ $? -ne 0 ]]; then
    return
  fi
  for branch in $(git branch | grep -v 'master\|main\|*'); do 
    if [[ "$REMOTES" != *"$branch"* ]]; then
      echo "git branch -D $branch"
      git branch -D $branch  
    fi
  done
}  
alias main=master
master() {
  default_branch=$(git remote show origin | grep 'HEAD branch' | cut -d' ' -f5)
  echo "git checkout $default_branch && git pull"
  git checkout $default_branch && git pull
}

# KUBERNETES STUFFs
alias k="command kubectl"
alias getNodePorts="kubectl get svc --all-namespaces -o go-template='{{range .items}}{{range.spec.ports}}{{if .nodePort}}{{.nodePort}} - {{.name}}{{\"\n\"}}{{end}}{{end}}{{end}}' "
alias show-contexts="kubectl config get-contexts"
alias test-pod="kubectl run samkirsch-test-shell --rm -i --tty --image ubuntu -- bash"
kube-ns() {
  command kubectl config set-context --current --namespace=$1
}
pod-pids() {
  echo """bash -c 'for pid in \$(ls /proc/ | grep -Eo "[0-9]+"); do f="/proc/\${pid}/cmdline"; if [[ -f \$f ]]; then echo "==========================="; echo \$pid; cat \$f; echo ""; echo ""; fi; done"""
}
kube-secret() {
  ns=$1
  secret=$2
  if [[ -z $ns ]] || [[ -z $secret ]]; then
    echo "This func takes positional args [NAMESPACE] and [SECRET_NAME]"
    echo "Example: kube-secret my-namespace secretName"
    return
  fi
  select value in $(kubectl -n $ns get secrets $secret -o yaml | yq -rc '.data | keys[]'); do 
    break;
  done
  echo "Getting secret [${secret}] value under .data.${value}"
  kubectl -n $ns get secrets $secret -o yaml | yq -r ".data.${value}" | base64 -d
}

# GCLOUD Stuff
# /Users/sxk3161/repos/github.com/samkirsch10/terminator/scripts/iTermProdWatcher.py &
export USE_GKE_GCLOUD_AUTH_PLUGIN=True
g-proj() {
  echo "$(cat ~/.config/gcloud/configurations/config_default | grep project | awk '{print $NF}')"
}

# alias gcp-switch="~/repos/github.com/samkirsch10/terminator/scripts/gprojects.sh"
alias cls="gcloud container clusters list"
alias force_drain="~/repos/github.com/samkirsch10/terminator/scripts/k8s_drain.sh"

alias force_drain="~/repos/github.com/samkirsch10/samkirsch10-terminator/scripts/k8s_drain.sh"
eval "$(direnv hook zsh)"

# export KUBECTL_FZF_PORT_FORWARD_LOCAL_PORT=8765

# source <(kubectl completion zsh)
# source ~/.kubectl_fzf.plugin.zsh


