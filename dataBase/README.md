## mysql


Este paquete permite la conexion a diversas bases de datos con su integracion del APM.


```go
import "github.com/architecture-it/go-platform/dataBase"
import "context"

func NewSQLRepository() SQLRepository {
	repository := &sqlRepository{
		dataRepository: database.NewDataRepository(os.Getenv("GORM_DRIVER"), os.Getenv("SQL_CONNECTION")),
	}
	return repository
}

func (repo *sqlRepository) ObtenerLocalidades(ctx context.Context) ([]models.Localidad, error) {
	var localidades []models.Localidad
	db := repo.dataRepository.GetDB(ctx)
	if db == nil {
		return localidades, errors.New("ocurrio un error con la conexion a la base")
	}
	err := db.Find(&localidades).Error
	return localidades, err
}

```
