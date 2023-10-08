# Setup
```bash
PYTHON_VERSION="python3.11"
sudo apt install -y software-properties-common
sudo add-apt-repository ppa:deadsnakes/ppa
sudo apt update
sudo apt install -y $PYTHON_VERSION
stat /usr/bin/python3 | grep -q 'symbolic link'
if [[ "$?" == "0" ]]; then
    ls -la /usr/bin/python3 | grep -q "$PYTHON_VERSION"
    if [[ "$?" != "0" ]]; then
        rm /usr/bin/python3
        ln -s /usr/bin/${PYTHON_VERSION} /usr/bin/python3
    fi
fi
curl https://bootstrap.pypa.io/get-pip.py -o get-pip.py
sudo python3 get-pip.py
pip install --user pipenv virtualenv
virtualenv ~/venv
```

## Install The Pretty Stuff
1) Install [Powerline fonts](https://github.com/powerline/fonts)
 - Run the install.sh
2) wget the [Terminator config](https://raw.githubusercontent.com/SamKirsch10/terminator/master/.config/terminator/config)
```bash
mkdir -p ~/.config/terminator
wget https://raw.githubusercontent.com/SamKirsch10/terminator/master/.config/terminator/config -o ~/.config/terminator/config
```
3) Quit/Reopen Terminator.
4) Install ZSH and set as default
```bash
sudo apt install zsh
wget https://github.com/robbyrussell/oh-my-zsh/raw/master/tools/install.sh -O - | zsh
```
5) ZSH will ask you some questions... we're doing nothing cause we'll use the files here. download / override the zsh files with the ones here
