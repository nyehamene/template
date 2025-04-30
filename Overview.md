# Proj

## Comment

## Literal

## Documentation

## Tag

## Package

```
p: package: package("package/name")
p :: package("package/name")
```

## Import

```
i: import: import("package/name")
i :: import("package/name")
```

## Constant

```
PI: Float: 3.14
PI :: 3.14
```

## Var

```
month_day: Int = 28
month_day := 28
```

## Function

```
add: fn(Int,Int)->Int: fn(a: Int, b: Int) -> Int {
  return a + b
  }
add :: fn(a: Int, b: Int) -> Int {
  return a + b
  }

// call a function
add(10, 5)
```

## Multiple return

## Function overload

```
add_int :: fn(a: Int, b: Int) -> Int { return a + b }
add_float :: fn(a: Float, b: Float) -> Float { return a + b }

// what type should an overload have?
add :: overload{ add_int; add_float }

// call an overloaded function

add(10, 5)     // forwards to add_int
add(10.0, 5.0) // forwards to add_float
```

## Split function

```
add_user->to_db :: fn(name: String, email: String)(db: sql.DB) {}

// call a split function

add_user("John", "john@mail.com")->to_db(my_db)
```

### Alternate syntax

```
add_user, to_database :: fn(user: User) (db: sql.DB) {}

add_user(User(name: "John")) to_database(sql.DB.open())

// or

add_user(...),
  to_database(...)

// or

add_user(
    ...
  ) to_database(
    ...
  )
```

## Record

```
Person: type: record{
  Name: String
  }
Person :: record{
  Name: String
  }

p :: Person{Name: "John" }
p.Name
```

## Composition

```
User :: record {
  using Person
  Email: String
  }

u :: User{ Name: "John", Email: "john@mail.com" }
u :: User{ Person{ Name: "John" }, Email: "john@mail.com" }
u.Name
u.Email
```

## Derived type

```
Person :: record{ Age: Int }
User: type: type(Person)
User :: type(Person)

assert(Person != User)
u :: User{}
u.Age

```

### Type conversion (explicity)

```
p :: Person{}
u :: User(p)
p :: Person(u)
```

## Type alias

```
User: type: Person
User :: Person

assert(User == Person)

p :: Person{}
// no type conversion needed
u: User: p
```

## Enum

An enum is a name given to a set of constant values.

```
Male :: 0
Female :: 1

Gender :: enum { Male; Female }

g1: Gender: Male
g2: Gender: Female
```

## Union

A union is a name given to a set of types

```
Guest :: record { Name: String }
Premium :: record { using Guest; Email: String }

User :: union { Guest; Premium }

u1: User: Guest{}
u2: User: Premium{}
```

## Method

```
User :: record{ Name: String }
to_string :: impl u: User fn() -> String { return u.Name }

u :: User{}
u.to_string()
```

### Method overload

```
Account :: record{ Balance: Float }
credit_int :: impl a: Account fn(amount: Int)  {}
credit_float :: impl a: Account fn(amount: Float) {}

credit :: impl Account overload{ credit_int; credit_float }

a :: Account{}
a.credit(10)   // forwards to Account.credit_int
a.credit(10.0) // forwards to Account.credit_float
```

### Method overload for union type

```
Dog :: record {}
Cat :: record {}

bark :: impl d: Dog fn() -> String {}
meow :: impl c: Cat fn() -> String {}

Animal :: union { Dog; Cat }

speak :: impl Animal overload { bark; meow }

animal : Animal : Dog{}
animal.speak() // forwards to Dog.bark

animal : Animal : Cat{}
animal.speak() // forwards to Cat.meow
```

### Method overload for enum type

```
True :: record{}
False :: record{}

if_true->else :: impl t: True fn() then: do() else: do() { goto then() }
if_false->else :: impl t: True fn() then: do() else: do() { goto else() }

if_true->else :: impl f: False fn() then: do() else: do() { goto else() }
if_false->else :: impl f: False fn() then: do() else: do() { goto then() }

// not a pubic API
TRUE :: True{}
FALSE :: False{}

Bool :: enum{ TRUE; FALSE }

// public API
true: Bool: TRUE
false: Bool: FALSE


if_true->else :: impl Bool overload {
  True.if_true->else;
  False.if_true->else
  }

if_false->else :: impl Bool overload {
  True.if_false->else;
  False.if_false->else
  }
```

## Generic (Parametric polymorphism)

### Polymorphic function

### Polymorphic function overload

### Polymorphic type

### Polymorphic union

### Polymorphic method overload for union
