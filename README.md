# btcon - Bluetooth Connect
Connect to bluetooth devices automatically when it is in range - a small fix to ubuntu's autoconnection bug.

The program polls the bluetooth device every 5 seconds and tries to connect to any of the listed devices.

## SetUp
1. clone the repo
2. edit main.go and mention ur devie MAC address in the `devicesList` slice
3. install it: `go install`
4. create a new service entry at `/etc/systemd/system`, say `btcon.service`

    ```sh
    nano /etc/systemd/system/btcon.service
    ```

    Enter the contents in the file:
    ```
    [Unit]
    Description=Bluetooth Auto Connection Service
    After=network.target
    StartLimitIntervalSec=0

    [Service]
    Type=simple
    Restart=always
    RestartSec=1
    User=YOUR_USER_NAME
    ExecStart=ABSOLUTE_GO_PATH/bin/btcon

    [Install]
    WantedBy=multi-user.target
    ```
5. start the service : `systemctl start btcon`
6. to run at startup enable it too: `systemctl enable btcon`

:|: any glitch? dont complain! :|: