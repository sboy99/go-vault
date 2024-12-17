package services

import "fmt"

func GetPortBoilerPlateContent(name string) string {
	return fmt.Sprintf(`
		export interface %s {
			// Add methods here
		}
	`, name)
}

func GetIndexBoilerPlateContent(fileName string) string {
	return fmt.Sprintf(`export * from './%s'`, fileName)
}

func GetControllerBoilerPlateContent(name string) string {
	return fmt.Sprintf(`	
		import { Controller, Get } from '@nestjs/common';

		@Controller()
		export class %s implements Controller {
		  
		}
	`, name)
}
