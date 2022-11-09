Package ints provides a fast Set-type for non-negative small
integers and a light weight UDecimal type with its associated
Context for non negative integer based decimal arithmetics.
Skimming through Set's API should suffice to learn what you
can do with provided Set-type.

To use the non-negative integer based decimal arithmetics a
Context instance is needed to create UDecimal values and
to operate on them.  UDec provides the ready to use default
context.

```go
    d1, err := ints.UDec.From.Str("19.5") // string to UDecimal
    if err != nil {
        panic(err)
    }

    d2, err := ints.UDec.From.Float(22.5) // float to UDecimal
    if err != nil {
        panic(err)
    }

    fmt.Println("result:", ints.UDec.MAdd(d1, d2).Str(ints.Dec))
    // prints "result: 42,00"
```


Note that UDecimal values are not aware of their number of fractionals,
i.e. positions after the decimal separator.  This information is kept in
the context.  Hence to get a string representation of a UDecimal value a
Context instance must be provided.

There are also 'Must'-versions if you feel save:

```go
    result := ints.UDec.MAdd(
        ints.UDec.From.MStr("2.5"),
        ints.UDec.From.MFloat(2.5),
    )
```

and the possibility to calculate directly with floats in a decimal
manner.

```go
    sum := ints.UDec.Float.MSum(2.5, 2.5)
```

I.e. provided floats are converted to UDecimal values before
the operation is performed and the operation returns an
Decimal.

A Context has two set of flags the arithmetic flags and the format
flags.  The former control the UDecimal conversion and their
arithmetics.  The later control the string representation of a
UDecimal value.  Once a Context instance is created its arithmetic
flags are immutable.  The UDec Context's arithmetic flags default to
DOT_SEPARATOR|SIX_FRACTIONALS.  I.e. the string to UDecimal
conversion expects a string with a dot decimal separator and the last
six positions of a converted Decimal represent its fractionals.  To
create a context with different arithmetic flags the UDec Context's
New method must be used

```go
    dec := ints.UDec.New(COMMA_SEPARATOR|FOUR_FRACTIONALS, DEFAULTS)
    d := dec.From.MStr("2,5")
```

The second argument sets the format flags which in above case are
simply a copy of UDec's format flags.

NOTE while I try to write idiomatic go code I value a convenient to
use (practical) API higher than the "make the zero type usable"
idiom.  I couldn't find a way to make the zero Context usable while
providing the API it has.  Hence you should never instantiate a
Context, Convert, Decimal or Floats type directly.  It is all
done for you by a single ints.UDec.New call.  The mentioned types are
defined as public anyway to provide their methods' documentations in
the go documentation server.

Format flags may be changed at any time

```go
    dec.SetFmt(DOT_SEPARATOR)
    fmt.Print(d.Str(dec)) // prints "2.50"
```

While Str truncates superfluous fractionals Rnd returns a "round to
even" string representation of a Decimal value.

Enjoy!
