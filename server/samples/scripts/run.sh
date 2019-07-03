virtualenv .env
source .env/bin/activate

python3 -m pip install regex # don't know why we need this. SoMaJo breaks otherwise
pip install -r requirements.txt
python main.py

deactivate