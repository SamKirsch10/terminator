# Setup
```bash
sudo apt install -y software-properties-common
sudo add-apt-repository ppa:deadsnakes/ppa
sudo apt update
sudo apt install -y python3.7 terminator dbus-x11
curl https://bootstrap.pypa.io/get-pip.py -o get-pip.py
sudo python3 get-pip.py
pip install --user pipenv virtualenv
virtualenv ~/venv
```



## WSL - To get terminator to open "natively"
![terminator](https://raw.githubusercontent.com/SamKirsch10/terminator/master/assets/window.png)

1) Install [VcXsrv](https://sourceforge.net/projects/vcxsrv) for X11 support


2) Make a new shortcut with the Target set to:
```
"C:\Program Files\VcXsrv\vcxsrv.exe" :0 -ac -terminate -lesspointer -multiwindow -clipboard -wgl -dpi auto 
```

3) Open `run.exe` and type `shell:startup`. Put the shortcut here.

4) Create a vbs script like below
```vbs
args = "-c" & " -l " & """cd ~; env DISPLAY=:0.0 terminator"""
WScript.CreateObject("Shell.Application").ShellExecute "bash", args, "", "open", 0
```

5) Create a shortcut with Target set to the script
```
C:\Windows\System32\wscript.exe D:\Documents\Scripts\Terminator.vbs
```

You can change the shortcut's icon with this one in this repo

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
