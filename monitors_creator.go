package main

// CreateMonitors создает модули мониторинга
func CreateMonitors() error {
	// создаем новый субъект и получаем его идентификатор в БД
	subjectID, err := CreateSubject()
	if err != nil {
		return err
	}

	// собираем данные для новых модулей мониторинга
	monitorModules, err := CollectionMonitorData(subjectID)
	if err != nil {
		return err
	}

	// перебираем карты с данными о модулях мониторинга и добавляем данные в БД
	for _, monitorModule := range *monitorModules {
		// добавляем в БД данные о новом модуле мониторинга
		newMonitor, err := createMonitor(monitorModule)
		if err != nil {
			return err
		}

		// собираем данные об использовании методов vk api
		methods, err := CollectionMethodData(newMonitor)
		if err != nil {
			return err
		}

		// добавляем собранные данные о методах в БД
		err = createMethods(methods)
		if err != nil {
			return err
		}

		// добавляем в таблицу соответствующего модуля мониторинга в БД новые данные
		err = createMonitorModules(monitorModule)
		if err != nil {
			return err
		}
	}

	return nil
}

// createMonitor выполняет алгоритмы создания монитора
func createMonitor(monitorModule MonitorModule) (*Monitor, error) {
	// добавляем новое поле в таблицу monitor
	var monitor Monitor
	monitor.Name = monitorModule.Name
	monitor.SubjectID = monitorModule.SubjectID
	err := InsertDBMonitor(monitor)
	if err != nil {
		return nil, err
	}

	// получаем из БД данные о новом мониторе
	newMonitor, err := SelectDBMonitor(monitorModule.Name,
		monitorModule.SubjectID)

	return newMonitor, nil
}

// createMethods выполняет алгоритмы добавления новых методов в БД
func createMethods(methods *[]Method) error {
	for _, method := range *methods {
		err := InsertDBMethod(method)
		if err != nil {
			return err
		}
	}
	return nil
}
