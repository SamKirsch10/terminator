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

Create a vbs script like below
```vbs
args = "-c" & " -l " & """cd ~; env DISPLAY=:0.0 terminator"""
WScript.CreateObject("Shell.Application").ShellExecute "bash", args, "", "open", 0
```

Create a shortcut with Target set to the script
```
C:\Windows\System32\wscript.exe D:\Documents\Scripts\Terminator.vbs
```

You can change the shortcut's icon with this one in this repo
