import * as kind from "@pulumi/kind-native";

const random = new kind.Random("my-random", { length: 24 });

export const output = random.result;