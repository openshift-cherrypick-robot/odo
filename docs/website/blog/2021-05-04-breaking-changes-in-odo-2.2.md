---
title: Breaking changes in odo 2.2
author: Girish Ramnani
author_url: https://github.com/girishramnani
author_image_url: https://avatars.githubusercontent.com/u/6551988?v=4
tags: ["release"]
slug: breaking-changes-in-odo-2.2
---

Breaking changes in odo 2.2

<!--truncate-->
### Breaking changes in odo 2.2

This document outlines the breaking changes that were introduced in odo 2.2. With the increased adoption of [Devfile](https://devfile.github.io/) we have started to reduce odo’s dependency on S2I(Source-to-Image). If you do not work with S2I, then you can stop reading here.

1. `odo create --s2i <component-type>` **will create a converted Devfile based component on the S2I images of that component type.**

   ```shell
   odo create --s2i nodejs
   ```

   Output -
   ```shell
   $ odo create nodejs --s2i
   Validation
    ✓  Validating component [424ms]
   Conversion
    ✓  Successfully generated devfile.yaml and env.yaml for provided S2I component

   Please use `odo push` command to create the component with source deployed
   ```
   
   The above command would generate a `devfile.yaml` which would be using the S2I images and variables that are part of the `nodejs`.

   This change will not break any existing S2I components. Although you are encouraged to convert them to devfile using `odo utils convert-to-devfile`.

2. **Currently devfile components do not support `--git` and `--binary` components hence we still use S2I component flow to create them.**

   ```shell
   $ odo create java --s2i --git ./build.war
   Validation
    ✓  Validating component [431ms]
   
   Please use `odo push` command to create the component with source deployed

   ```
   Observe that there was no conversion done here.


3.  **`odo env set DebugPort` won't work with converted devfile components, you would need to use `odo config set --env DEBUG_PORT` instead.**

   Currently, the `wildfly` and `dotnet` component types do not work when converted. We have an issue open for this - <https://github.com/redhat-developer/odo/issues/4623>

### Known bugs and limitations
- https://github.com/redhat-developer/odo/issues/4623
- https://github.com/redhat-developer/odo/issues/4615
- https://github.com/redhat-developer/odo/issues/4594
- https://github.com/redhat-developer/odo/issues/4593


### Frequently asked questions
1. Why does odo fail create to URL using `odo url create` for a component created using `odo create --s2i` even though `odo url create` is allowed for devfile?

   * It won’t fail in the sense that if you tried the conventional s2i approach and try to create `odo url create` it would fail with url for 8080 port already present as there would already be one for you. Refer - https://github.com/redhat-developer/odo/issues/4621

2. How to understand the status of the debug?

   * `odo env set DebugPort` won't work, instead you would need to use `odo config set --env DEBUG_PORT` - this is because the s2i to devfile converted devfiles don't have a debug type command defined in them. We would fix this too.

3. If every component will now be reported as devfile component, then what about existing S2I components?

   * They should work as is. Check `odo list` for a simpler check. Somethings might still break since it is quite complex to make things work across the board.


4. If oc based checks in tests are not going to work, is there an alternative odo support around it ?
   
   * It would’t be right to say that you cannot use `oc` based checks, but they would break because now the s2i components are being converted to devfile, odo would generate a Kubernetes `Deployment` but the `oc` would try to find a `DeploymentConfig` on the cluster.
