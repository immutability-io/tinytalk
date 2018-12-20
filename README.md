# TinyTalk

You need to run `tinytalk` from the directory where you have your `tinytalk.creds` file. Obviously, this could be a parameter, but it isn't so deal with it.

## tinytalk.creds

I will give you a file encrypted with your Keybase identity. Suppose your keybase identity is `keybase:mcboaterson`. I will give you a file named `mcboaterson.creds.enc`. You will need to decrypt that file into a file named `tinytalk.creds`.

For example:

```
$  keybase pgp decrypt -i mcboaterson.creds.enc > /path/somewhere/tinytalk.creds
```

## tinytalk executable

I've built a linux and darwin (mac) version of the `tinytalk` program. You can pull it from the releases directory. You can put it in your PATH or you can put it in the same directory where `tinytalk.creds` is. You can download from this repo's releases or you can use `go get` to build the program.

## Listen for messages

To listen for messages, change to the directory where `tinytalk.creds` reside and run:

```
$ tinytalk listen &
```

## Send an unencrypted message

To send plain text messages, change to the directory where `tinytalk.creds` reside and run:

```
$ tinytalk say -m "this is what i want to say"
```

If you want to identify yourself as the sender, use this form of the command:

```
$ tinytalk say -s "boaty" -m "this is what i want to say"
```


## Send an unencrypted message

To send a private messages, change to the directory where `tinytalk.creds` reside and run:

```
$ tinytalk say -m "this is what i want to say" -k "keybase:mcboaterson"
```

In this case, the `keybase` identity is the target of the secret message. If you want to identify yourself as the sender, use this form of the command:

```
$ tinytalk say -s "boaty" -m "this is what i want to say" -k "keybase:mcboaterson"
```
