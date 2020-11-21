func WithTransactionExample() {
	ctx := context.Background()
	For a replica set, include the replica set name and a seedlist of the members in the URI string; e.g.
	uri := "mongodb://mongodb1.localhost:27001,mongodb2.localhost:27002,mongodb3.localhost:27003/?replicaSet=auten"
	For a sharded cluster, connect to the mongos instances; e.g.
	uri := "mongodb://mongos1.localhost:27001,mongos2.localhost:27002,mongos2.localhost:27003"
	var uri string

	clientOpts := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		panic(err)
	}
	defer func() { _ = client.Disconnect(ctx) }()

	Prereq: Create collections.
	wcMajority := writeconcern.New(writeconcern.WMajority(), writeconcern.WTimeout(1*time.Second))
	wcMajorityCollectionOpts := options.Collection().SetWriteConcern(wcMajority)
	fooColl := client.Database("med").Collection("foo", wcMajorityCollectionOpts)
	barColl := client.Database("med").Collection("bar", wcMajorityCollectionOpts)

	Step 1: Define the callback that specifies the sequence of operations to perform inside the transaction.
	callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
		// Important: You must pass sessCtx as the Context parameter to the operations for them to be executed in the
		// transaction.
		if _, err := fooColl.InsertOne(sessCtx, bson.D{{"abc", 1}}); err != nil {
			return nil, err
		}
		if _, err := barColl.InsertOne(sessCtx, bson.D{{"xyz", 999}}); err != nil {
			return nil, err
		}

		return nil, nil
	}

	Step 2: Start a session and run the callback using WithTransaction.
	session, err := client.StartSession()
	if err != nil {
		panic(err)
	}
	defer session.EndSession(ctx)

	result, err := session.WithTransaction(ctx, callback)
	if err != nil {
		panic(err)
	}
	fmt.Printf("result: %v\n", result)
}