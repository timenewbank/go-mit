## Go M.I.T Chain

Official golang implementation of the M.I.T Chain protocol.

Automated builds are available for stable releases and the unstable master branch.
Binary archives will be published on official website soon..

## Building the source

For prerequisites and detailed build instructions please read the
Installation Instructions

Building mit requires both a Go (version 1.9 or later) and a C compiler.
You can install them using your favourite package manager.
https://golang.org/doc/install

Once the dependencies are installed, run

    cd go-mit
    make mit

or, to build the full suite of utilities:

    make all

## Running Mit

Going through all the possible command line flags is out of scope here,but we've
enumerated a few common parameter combos to get you up to speed quickly on how you can run your
own Mit instance.

### Full node on the main Mit network

By far the most common scenario is people wanting to simply interact with the Mit network:
create accounts; transfer funds; deploy and interact with contracts. For this particular use-case
the user doesn't care about years-old historical data, so we can full-sync quickly to the current
state of the network. To do so:

```
$ mit console
```
#### Running a miner

Currently our main network is based on POW+POS.

This command will:

 * Start mit in full sync mode (default), causing it to
   download more data in exchange for avoiding processing the entire history of the Mit network,
   which is very CPU intensive.
 * Start up Mit's built-in interactive JavaScript console,
   (via the trailing `console` subcommand) through which you can invoke all official `web3` methods
   as well as Mit's own management APIs.
   This too is optional and if you leave it out you can always attach to an already running Mit instance
   with `mit attach`.

## Running mit testnet

Going through all the possible command line flags is out of scope here , but we've
enumerated a few common parameter combos to get you up to speed quickly on how you can run your
own Mit instance.

### Full node on the Mit testnet

By far the most common scenario is people wanting to simply interact with the Mit network:
create accounts; transfer funds; deploy and interact with contracts. For this particular use-case
the user doesn't care about years-old historical data, so we can fast-sync quickly to the current
state of the network. To do so:

```
$ mit console
```

This command will:

 * Start mit in fast sync mode (default, can be changed with the `--syncmode` flag), causing it to
   download more data in exchange for avoiding processing the entire history of the Mit network,
   which is very CPU intensive.
 * Start up Mit's built-in interactive JavaScript console,
   (via the trailing `console` subcommand) through which you can invoke all official `web3` methods
   as well as Mit's own management APIs.
   This too is optional and if you leave it out you can always attach to an already running Mit instance
   with `mit attach`.

### Full node on the Mit test network

Transitioning towards developers, if you'd like to play around with creating Mit time contracts, you
almost certainly would like to do that without any real money involved until you get the hang of the
entire system.

```
$ mit --testnet console
```

Specifying the `--testnet` flag however will reconfigure your Mit instance a bit:

 * Instead of using the default data directory (`~/.timenewbank` on Linux for example), Mit will nest
   itself one level deeper into a `testnet` subfolder (`~/.timenewbank/testnet` on Linux). Note, on OSX
   and Linux this also means that attaching to a running testnet node requires the use of a custom
   endpoint since `mit attach` will try to attach to a production node endpoint by default. E.g.
   `mit attach <datadir>/testnet/mit.ipc`. Windows users are not affected by this.


*Note: Now we have provided the testnet for you. The mainnet will be launched in the near future. Please make sure that you always using two separate accounts for play-money and real-money respectively.*

#### Starting up your member nodes

With the bootnode operational and externally reachable (you can try `telnet <ip> <port>` to ensure
it's indeed reachable), start every subsequent Mit node pointed to the bootnode for peer discovery
via the `--bootnodes` flag.

```
$ mit --datadir=path/to/custom/data/folder --bootnodes=<bootnode-enode-url-from-above>
```

*Note: Since your network will be completely cut off from the main and test networks, you'll also
need to configure a miner to process transactions and create new blocks for you.*

#### Running a miner

Currently our testnest is based on POA(Proof-of-Authority). So the miners can only mining after getting permission from us.

## Contribution

Thank you for considering to help out with the source code! We welcome contributions from
anyone on the internet, and are grateful for even the smallest of fixes!

If you'd like to contribute to go-mit, please fork, fix, commit and send a pull request
for the maintainers to review and merge into the main code base.

to ensure those changes are in line with the general philosophy of the project and/or get some
early feedback which can make both your efforts much lighter as well as our review and merge
procedures quick and simple.

## License

The go-mit library (i.e. all code outside of the `cmd` directory) is licensed under the
[GNU Lesser General Public License v3.0](https://www.gnu.org/licenses/lgpl-3.0.en.html), also
included in our repository in the `COPYING.LESSER` file.

The go-mit binaries (i.e. all code inside of the `cmd` directory) is licensed under the
[GNU General Public License v3.0](https://www.gnu.org/licenses/gpl-3.0.en.html), also included
in our repository in the `COPYING` file.
