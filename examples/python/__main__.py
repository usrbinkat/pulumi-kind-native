import pulumi
import pulumi_kind as kind

my_arbitrary_kind_cluster_resource_name = kind.Kind("myArbitraryKindClusterResourceName", name="test")
