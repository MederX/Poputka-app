package handlers

import "cloud.google.com/go/firestore"

// firestoreMergeAll returns the MergeAll option for Set calls.
func firestoreMergeAll() firestore.SetOption {
	return firestore.MergeAll
}
