# Language reference

## Explicit type inference

(Hopes that makes sense)

A use can request explicit type inference using `type` keyword in
the type position in a declaration. For example,

### Explicit type inference for custom types

```
User : type : type(Person)
```
In the above declaration the type of `User` is infered from the expression
on the RHS of the declaration.
In this case the type is infered to be `type(Person)`.

### Explicit type inference for record type

```
Person : type : record{
  name: String
  email: String
  }
```
In the above declaration the type of `Person` is infered from the RHS
expression. In this case the infered type is `record{name: String; email: String}`
coincidentally (almost) matching the record declaration literal on the RHS. 

## Complete declaration syntax

Below describes the complete declaration without omitting any optional
parts and not including explicit type inference.

### Declaring a custom type based on another type

The complete declaration of a custom type is as follows using the `User`
type declared in #[Explicit type inference]:

```
User : type(Person) : type(Person);
```

Here, starting from the RHS, `type(Person)` as an expression that resolves
to a type for the existing type `Person`. The type the expression is `type(Person)`,
then the resolved type is stored in the identifier `User`.

### Declaring a record type

The complete syntax for declaring a record type is very verbose and should
be avoided. Explicit or implicity type inference should be prefered.
The syntax is as described below:

```
Person : record{name: String; email: String} : record {
  name: String
  email: String
  }
```

Here, on the RHS is the record declaration literal declaring all the fields
in the record with their associated types. Next, the middle section is the
the type for the record. In this case it repeats exactly the same information
as in the record declaration literal. However, this is not the case all
situations. A type type declaration cannot contain tags, and documentations.
On the LHS `Person` is an identifier pointing to the type resolved from the
RHS expression.

### Declaring a templ

The complete syntax for declaring a template is as follows:

```
User : templ : templ(u: User) {
  <p
    Hello, (u.name)!
    Your email is (u.email).
    />
  }
```

Starting from the RHS, is a templ declaration literal declaring all the
elements, components and texts in the template. The literal beging wit
the keyword `templ`, then a parameter list of only one parameter declaring
a variable whose type must match the type for which the template is defined
for. Last, a block of template element/text/component.

The `templ` keyword is overload just as with `package`, `import`, and used
to declare the type of the template. On the LHS `User` is an identifer to
stores/points to the template resolved from the RHS expression.

### Templ declaration and explicit type inference

A template declaration does not declare a type therefore the keyword `type`
cannot be used for explicit type inference. In fact the templ declaration
syntax does not include explicit type inference.

However, the type of the parameter list can be explicitly infered since it
always matches the type on the LHS.

```
User :: templ(u: type) {
  <p
    Hello, (u.name)!
    Your email is (u.email).
    />
  }
```

### Templ declaration and implicit type inference

A template declaration type can be infered implicitly as follows:

```
User :: templ(u: User) {
  <p
    Hello, (u.name)!
    Your email is (u.email).
    />
  }
```

#### Special syntax for declaring template

There is an additional syntax for declaring template. And this is the
recommended syntax.

```
User :: (u: type) {
  <p
    Hello, (u.name)!
    Your email is (u.email).
    />
  }
```
