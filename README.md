# Cryptocurrency Trading Bot in Go

This repository is part of a YouTube series on creating a cryptocurrency
trading bot in Go. To follow this series, please visit the [Youtube channel](https://www.youtube.com/channel/UCWQaM7SpSECp9FELz-cHzuQ)

## Requirements

### Git (Required)

In order to checkout this code, first install Git and clone the code using
the tools above.

### Go (Required)

In order to get started, please install [Go](https://go.dev/dl/) version 1.19
or above.

### Make (Required)

#### Windows

Please refer to [this guide](https://stackoverflow.com/questions/32127524/how-to-install-and-use-make-in-windows) in order to install on Windows.

#### Linux

You'll likely know how to do this, but refer to your package manager for
instructions. If you're on ubuntu, it'll probably be:

```
$ sudo apt-get -y install make
```

#### macOS

I've not used macOS in a number of years, but you should be able to use
Make once you've installed [Xcode developer tools](https://www.freecodecamp.org/news/install-xcode-command-line-tools/). 

### golanglint-ci (Optional)

To be able to lint the code, you'll need to install the golanglint-cli package
as per the [instructions](https://golangci-lint.run/usage/install/#local-installation).


## Running the code

You should watch the series before attempting to run this code. It's unfinished
in it's current form and will not work without you making some adjustments. 

In case you need a reminder, here are the following supported commands

### Help

All the commands are viewable by running 

```
make help
```

### Building

```
make build
```

### Testing

```
make test
```

### Running

```
make run
```

### Linting

```
make lint
```

## FAQs

### Will this make me rich from trading?

No, this code will not make you rich. In order for that to happen, you'll need
to be an expert at trading, and even then you're up against other experts. 

This code will help you to understand how to create trading bots and to become
a software developer. This will give you the necessary tools to build your 
own trading strategies.

### Can you write a strategy that will make me rich?

No, whilst I know how to write software, I am a terrible trader. I have my own
strategies when it comes to investments and I stick to them (for better or worse). 
