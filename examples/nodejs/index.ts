import * as pulumi from "@pulumi/pulumi";
import * as kind from "@pulumi/kind";

const myArbitraryKindClusterResourceName = new kind.Kind("myArbitraryKindClusterResourceName", {name: "test"});
