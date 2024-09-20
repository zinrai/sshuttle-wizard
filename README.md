# sshuttle-wizard

sshuttle-wizard is an interactive CLI tool designed to simplify the configuration and execution of sshuttle commands. It automatically detects private subnets on a remote host, allows users to easily select them, and then generates and optionally executes the sshuttle command.

## Features

- Automatic detection of private subnets via SSH connection to the remote host
- Support for selecting multiple subnets
- Automatic generation of sshuttle commands
- Option to directly execute the generated command

## Prerequisites

- sshuttle (installed and available in your system PATH)
- SSH access to the remote host

## Installation

2. Build the tool:

   ```
   $ go build
   ```

## Usage

1. Run sshuttle-wizard in your terminal:

   ```
   ./sshuttle-wizard
   ```

2. Follow the prompts to enter:
   - The remote host (e.g., user@example.com)
   - Subnets to route (select from the automatically detected list)
   - Additional sshuttle options (if needed)

3. Review the generated command and choose whether to execute it.

## Example

```
$ ./sshuttle-wizard
sshuttle-wizard
---------------
Welcome to the sshuttle command builder and executor wizard!
Enter remote host (e.g. user@example.com): user@example.com

Detected private subnets:
1. 10.0.0.0/24
2. 172.16.0.0/16
3. 192.168.0.0/24
Enter the numbers of the subnets to route (comma-separated, or press Enter to finish): 1,3
Enter additional options (e.g. -v for verbose): -v

Prepared sshuttle command:
sshuttle -v 10.0.0.0/24 192.168.0.0/24 -r user@example.com

Do you want to execute this command? (y/n): y
Executing sshuttle command...
```

## License

This project is licensed under the MIT License - see the [LICENSE](https://opensource.org/license/mit) for details.
