package container

type Service struct{}

type ServiceInterface interface {
	ListContainer() ([]ContainerDTO, error)
	StartContainer(containerId string) error
	StopContainer(containerId string) error
	PublishContainer() error
}
