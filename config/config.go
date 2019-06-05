package config

const (
	ClientId        = "a1ce2967-a534-417d-bf12-c86571e4eefa"
	ClientSecret    = "a3f52f6c417f24bfdf583ed884c5d0cb489320c58222b061298e4a2d41a1bbd7"
	PIN 						= "457965"
	PinToken        = "0t4EG7tJerZYds7N9QS0mlRPCYsEVTQBe9iD1zNBCFN/XO7XEB87ypsCDWfRmDiZ7izzB/nokuMJEu6RJShMHCdIwYISU9xckA/8hIsRVydvoP14G/9kRidMHl/3RPLDMK6U2yCefo2BH0kQdbcRDxpiddqrMc4fYmZo6UddU/A="
	SessionId       = "26ed1f52-a3b4-4cc3-840f-469d3f19b10b"
  PrivateKey      = `-----BEGIN RSA PRIVATE KEY-----
MIICXAIBAAKBgQDaSPE8Cu18xzr8MOcgJx8tQnRdlS7c6JVs23497IGdIybIUYmZ
8zvgrFozpGjQYz2ayRDMWUQd/wm7e0Tf7n4bVCmQfkk72usAHX6pNA4HUeTeTmDT
sZQKdVx0K84Y3u512cAi5artnUjIsFRPP/LhAX0ujdgNMWIcHrMRh77s1wIDAQAB
AoGAVPW3Dwuhy8MvriDKlLUlaVRIPnRmPQ05u5ji1e9Ls4GPAsDZsdX+JEBxC1Ce
ix1VSP2hUCgeXx55B0O/VvlYk0pfogrxDgOw2dP04uboMG7tSE4TZK8J9zFPUrE0
wizFmbkgV2OEw33r00FqEhr0KnB9kXOzB5BvKN/FVyXui+ECQQDz1x3hOypW2kM9
uOqjQyg55VDkkXVZ8RgOmVd24MfkDjRauj1oGgLUWvINzhmXN5m84IhlOz1hgEuO
enHOpMmDAkEA5SuVeRhBZofUoaRbFxWL4jAN6+uuxFxZ0gCc9l4gwFkQp0RbEw/S
tiX9Cl06JR2oc2FBlaO5Vi1u8XfxOSUzHQJBANijfKaJHFrB3A/QZJbcqbaWaEJK
gYqBSzBdSHoTx0R04krhQIFm6rCkhH2DaPUSrwJCMqxN74DarUZOvyIrAeUCQH2F
ecFx/6BhFZ3Tn/Ds5ElneLiXxonW63uSymZG+DlijzSOxDOUnx0VgZuDpK1fqTxJ
MNr9ai5BhFrOD1n1fiECQBafDxsfFQv3w6j5/2PL54DhddGo50FzGxYR1LlttdVI
Q04EytqK7grDDS9PsfeXqdUo0D3NMSJ0BYs/kDsqGSc=
-----END RSA PRIVATE KEY-----`
)
const (
	SqlitePath = "./payment.db"
	//check the pending order second
	CheckPendingOrderInterval = 6
	//minutes
	OrderExpired              = 12
)

var Assets = map[string]bool{
	"815b0b1a-2764-3736-8faa-42d694fa620a": true,
	"6cfe566e-4aad-470b-8c9a-2fd35b49c68d": true,
	"965e5c6e-434c-3fa9-b780-c50f43cd955c": true,
}
