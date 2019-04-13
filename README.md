# Setup
```bash
sudo apt install software-properties-common
sudo add-apt-repository ppa:deadsnakes/ppa
sudo apt update
sudo apt install python3.7 terminator dbus-x11
curl https://bootstrap.pypa.io/get-pip.py -o get-pip.py
sudo python3 get-pip.py
pip install --user pipenv virtualenv
virtualenv ~/venv
```



## WSL - To get terminator to open "natively"

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
