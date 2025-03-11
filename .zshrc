# If you come from bash you might have to change your $PATH.
export PATH=~/src/terminator/scripts:/usr/local/opt/make/libexec/gnubin:/opt/homebrew/bin:$HOME/go/bin:$HOME/bin:/usr/local/bin:/usr/local/go/bin:$PATH

export PATH=$PATH:/opt/homebrew/Cellar/libpq/17.2/bin

# Apache Maven Environment Variables
# MAVEN_HOME for Maven 1 - M2_HOME for Maven 2
export M2_HOME=/usr/local/src/apache-maven
export PATH=${M2_HOME}/bin:${PATH}

# Path to your oh-my-zsh installation.
export ZSH="$HOME/.oh-my-zsh"

eval $(brew shellenv)

export PATH=$HOME/src/fsdev/tools/node/bin/:$HOME/src/fsdev/tools/python3/bin:$HOME/src/fsdev/tools/helm/:$HOME/src/fsdev/tools/terraform:$HOME/src/mn/projects/fullstory/bin/:~/go/bin:$PATH
export PATH=$PATH:$HOME/.local/bin

export GOROOT=/Users/samkirsch/go

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

# export ZSH_DISABLE_COMPFIX="true"
source $ZSH/oh-my-zsh.sh

# export ANSIBLE_VAULT_PASSWORD_FILE=/home/skirsch/v

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
alias dive="docker run -ti --rm  -v /var/run/docker.sock:/var/run/docker.sock wagoodman/dive"
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
alias git-sha="git rev-parse HEAD"
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
rebase() {
  default_branch=$(git remote show origin | grep 'HEAD branch' | cut -d' ' -f5)
  current_branch=$(git rev-parse --abbrev-ref HEAD)
  if [[ "$default_branch" == "$current_branch" ]]; then
    echo "Can't rebase the $default_branch branch!"
    return
  fi
  master
  git checkout $current_branch
  git rebase $default_branch
}
green() {
  git checkout green && git pull
}
stash() {
  MSG="$1"
  if [[ -z "$MSG" ]]; then
    echo "$0 expects a message!"
    return
  fi
  echo "git stash save \"$MSG\""
  git stash save "$MSG"
}

# KUBERNETES STUFFs
export KUBECTL_EXTERNAL_DIFF="dyff between --omit-header --set-exit-code"
alias k="command kubectl"
alias getNodePorts="kubectl get svc --all-namespaces -o go-template='{{range .items}}{{range.spec.ports}}{{if .nodePort}}{{.nodePort}} - {{.name}}{{\"\n\"}}{{end}}{{end}}{{end}}' "
alias show-contexts="kubectl config get-contexts"
alias test-pod="kubectl run samkirsch-test-shell --rm -i --tty --image us-central1-docker.pkg.dev/fs-ops/fullstory-images/debugger -- bash"
alias klogs="command kubectl logs -f"
kube-ns() {
  command kubectl config set-context --current --namespace=$1
}
pod-pids() {
  echo """bash -c 'for pid in \$(ls /proc/ | grep -Eo "[0-9]+"); do f="/proc/\${pid}/cmdline"; if [[ -f \$f ]]; then echo "==========================="; echo \$pid; cat \$f; echo ""; echo ""; fi; done"""
}
kube-secret() {
  ns=$1
  secret=$2
  if [[ -z $secret ]]; then
    secret="$ns"
    ns="default"
    echo "using namespace 'default'."
  fi

  if [[ -z $ns ]] || [[ -z $secret ]]; then
    echo "This func takes positional args [NAMESPACE] and [SECRET_NAME]"
    echo "Example: kube-secret my-namespace secretName"
    return
  fi
  select value in $(kubectl -n $ns get secrets $secret -o yaml | yq -rc '.data | keys[]'); do 
    break;
  done
  echo "Getting secret [${secret}] value under .data.${value}"
  kubectl -n $ns get secrets $secret -o yaml | yq -r '.data["'${value}'"]' | base64 -d
}
kube-cronjob() {
  name=$1
  action=$2
  suspend=""
  if [[ "$action" == "suspend" ]]; then
    suspend="true"
  elif [[ "$action" == "resume" ]]; then
    suspend="false"
  else 
    echo "kube-cronjob expects arguments with an action of suspend/resume"
    echo "example: kube-cronjob samsJob suspend"
    echo -e "\nYou can pass additional arguments to the resulting kubectl cmd after the action argument (like namespace etc)"
    echo "example: kube-cronjob samsJob suspend -n someNamespace"
    return 
  fi
  shift;
  shift;
  kubectl patch cronjob $name -p "{\"spec\" : {\"suspend\" : $suspend }}" "$@"
}
prom-op() {
  pod=$(kubectl -n operators get pods | grep prometheus-operator | awk '{print $1}')
  kubectl -n operators logs -f $pod
}

# GCLOUD Stuff
g-proj() {
  echo "$(cat ~/.config/gcloud/configurations/config_default | grep project | awk '{print $NF}')"
}

# alias gcp-switch="~/repos/github.com/samkirsch10/terminator/scripts/gprojects.sh"
alias cls="gcloud container clusters list"
alias force_drain="~/repos/github.com/samkirsch10/terminator/scripts/k8s_drain.sh"


alias fs-playpen="fsk8s fs-playpen"
alias fs-ops="fsk8s fs-ops"
alias fs-staging="fsk8s fs-staging"
alias fullstoryapp="fsk8s fullstoryapp"

export SKIP_FS_PS1=1
export FS_SKIP_CD=1
source /Users/samkirsch/.fsprofile
eval "$(direnv hook zsh)"


alias ssh="/Users/samkirsch/src/mn/projects/fullstory/tools/util/fsssh.sh"
