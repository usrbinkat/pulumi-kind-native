using System.Collections.Generic;
using System.Linq;
using Pulumi;
using Kind = Pulumi.Kind;

return await Deployment.RunAsync(() => 
{
    var myArbitraryKindClusterResourceName = new Kind.Kind("myArbitraryKindClusterResourceName", new()
    {
        Name = "test",
    });

});

