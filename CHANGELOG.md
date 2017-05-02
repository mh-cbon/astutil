# Changelog - astutil

### 0.0.4

__Changes__

- MethodParams: handle func type

__Contributors__

- mh-cbon

Released by mh-cbon, Tue 02 May 2017 -
[see the diff](https://github.com/mh-cbon/astutil/compare/0.0.4-beta...0.0.4#diff)
______________

### 0.0.4-beta

__Changes__

- MethodParams: it now handles starExpr, arrayExpr and ellipsisExpr.
- added MethodParamsToProps: useful to create a struct/map.

__Contributors__

- mh-cbon

Released by mh-cbon, Sun 30 Apr 2017 -
[see the diff](https://github.com/mh-cbon/astutil/compare/0.0.3...0.0.4-beta#diff)
______________

### 0.0.3

__Changes__

- various fixes in
  - MethodReturnPointer
  - MethodHasEllipse
  - MethodParamNamesInvokation







__Contributors__

- mh-cbon

Released by mh-cbon, Sat 29 Apr 2017 -
[see the diff](https://github.com/mh-cbon/astutil/compare/0.0.3-beta4...0.0.3#diff)
______________

### 0.0.3-beta4

__Changes__

- fix MethodReturnPointer

__Contributors__

- mh-cbon

Released by mh-cbon, Sat 29 Apr 2017 -
[see the diff](https://github.com/mh-cbon/astutil/compare/0.0.3-beta3...0.0.3-beta4#diff)
______________

### 0.0.3-beta3

__Changes__

- fix function MethodHasEllipse when the func does not have any parameters

__Contributors__

- mh-cbon

Released by mh-cbon, Sat 29 Apr 2017 -
[see the diff](https://github.com/mh-cbon/astutil/compare/0.0.3-beta2...0.0.3-beta3#diff)
______________

### 0.0.3-beta2

__Changes__

- fix function MethodParamNamesInvokation

__Contributors__

- mh-cbon

Released by mh-cbon, Sat 29 Apr 2017 -
[see the diff](https://github.com/mh-cbon/astutil/compare/0.0.3-beta1...0.0.3-beta2#diff)
______________

### 0.0.3-beta1

__Changes__

- Add new function
  - __MethodParamNamesInvokation__(*ast.FuncDecl, withEllipse) bool: return `s...` with `func(s ...string){}`














__Contributors__

- mh-cbon

Released by mh-cbon, Sat 29 Apr 2017 -
[see the diff](https://github.com/mh-cbon/astutil/compare/0.0.3-beta...0.0.3-beta1#diff)
______________

### 0.0.3-beta

__Changes__

- Add new function
  - __MethodHasEllipse__(*ast.FuncDecl) bool: return true if the last params uses ellipse.
- Initialize tests

__Contributors__

- mh-cbon

Released by mh-cbon, Sat 29 Apr 2017 -
[see the diff](https://github.com/mh-cbon/astutil/compare/0.0.2...0.0.3-beta#diff)
______________

### 0.0.2

__Changes__

- Add new functions
  - __IsBasic__(string) bool: to konw if t is string/int...
  - __GetPointedType__(string) string: Given `T|*T`, returns `*T`
  - __GetUnpointedType__(string) string: Given `T|*T`, returns `T`
  - __IsAPointedType__(string) bool: Given `*T`, returns `true`


















__Contributors__

- mh-cbon

Released by mh-cbon, Sat 29 Apr 2017 -
[see the diff](https://github.com/mh-cbon/astutil/compare/0.0.1...0.0.2#diff)
______________

### 0.0.1

__Changes__

- Initialize the project.

__Contributors__

- mh-cbon

Released by mh-cbon, Sat 29 Apr 2017 -
[see the diff](https://github.com/mh-cbon/astutil/compare/128ad89fb09c52948212c066b986977f43a2c8c1...0.0.1#diff)
______________


