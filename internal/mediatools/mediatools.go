package mediatools

import (
	"fyne.io/fyne/v2"
	"github.com/Developpeur-du-dimanche/MediaTools/internal/view"
	"github.com/Developpeur-du-dimanche/MediaTools/pkg/application"
)

func Run() {
	application := application.NewApp(&fyne.AppMetadata{
		ID:   "com.github.developpeur-du-dimanche.mediatools",
		Name: "MediaTools",
	})

	homeView := view.NewHomeView(application)

	homeView.LogGui()

}
