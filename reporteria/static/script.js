document.addEventListener("htmx:afterSwap", (event) => {
    if (event.target.id === "main-content") {
        // Lógica de inicialización de tu contenido , cunado el main-content sufre cambios
        console.log("Nuevo contenido cargado en #main-content");

        // activa la funcionalidad de los filtros
                                  // Variables
                            const filtersCard = document.getElementById('filtersCard');
                            const closeFiltersBtn = document.getElementById('closeFiltersBtn');
                            const showFiltersBtn = document.getElementById('showFiltersBtn');
                            
                            // Filtros manualmente definidos
                            const customSelects = [
                                { id: 'Sexo', selectedText: 'selectedTextSexo' },
                                { id: 'Pueblos Originarios', selectedText: 'selectedTextpueblosOriginarios' },
                                { id: 'Nacionalidad', selectedText: 'selectedTextNacionalidad' },      
                                { id: 'Zona Geografica', selectedText: 'selectedTextZonaGeografica' },
                                { id: 'Dependencia', selectedText: 'selectedTextDependencia' },
                                { id: 'Ivm', selectedText: 'selectedTextIvm' },
                                { id: 'Estado Nutricional', selectedText: 'selectedTextEstadoNutricional' },
                                { id: 'Talla', selectedText: 'selectedTextTalla' },
                                { id: 'Discapacidad', selectedText: 'selectedTextDiscapacidad' },
                                { id: 'Grupos Vulnerables', selectedText: 'selectedTextGruposVulnerables' }

                            ];
                            
                            // Mostrar/Ocultar opciones al hacer clic en cada select
                            customSelects.forEach((customSelect) => {
                                const selectElement = document.getElementById(customSelect.id);
                                const optionsContainer = selectElement.querySelector('.options-container');
                                const selectedTextElement = document.getElementById(customSelect.selectedText);
                                const defaultText = selectedTextElement.textContent; // Guardamos el texto predeterminado
                            
                                selectElement.addEventListener('click', (e) => {
                                    e.stopPropagation();
                                    selectElement.classList.toggle('open');
                                });
                            
                                // Escuchar eventos de clic en los checkboxes
                                optionsContainer.addEventListener('change', () => {
                                    const selectedOptions = Array.from(
                                        optionsContainer.querySelectorAll('.category-checkbox:checked')
                                    ).map(option => option.value);
                            
                                    if (selectedOptions.length > 0) {
                                        selectedTextElement.textContent = selectedOptions.join(', ');
                                    } else {
                                        // Restablecer el texto predeterminado si no hay seleccionados
                                        selectedTextElement.textContent = defaultText;
                                    }
                                });
                            });
                            
                            
                            // Cerrar las opciones si se hace clic fuera
                            document.addEventListener('click', (e) => {
                                customSelects.forEach(customSelect => {
                                    const selectElement = document.getElementById(customSelect.id);
                                    if (!selectElement.contains(e.target)) {
                                        selectElement.classList.remove('open');
                                    }
                                });
                            });
                            
                            // Manejador del botón "Cerrar"
                            closeFiltersBtn.addEventListener('click', () => {
                                filtersCard.classList.add('d-none');
                            });
                            
                            // Manejador del botón "Mostrar Filtros"
                            showFiltersBtn.addEventListener('click', () => {
                                filtersCard.classList.remove('d-none');
                            });
                            

                            // Manejador de envío de formulario
                            document.getElementById('filtersForm').addEventListener('submit', (e) => {
                                e.preventDefault();
                                
                                const filters = {};
                            


                                    // Obtener texto de las opciones seleccionadas de otros selects en la página
                                    const additionalSelects = document.querySelectorAll('select'); // Cambia el selector si tienes un conjunto específico
                                    additionalSelects.forEach((select) => {
                                        const selectedOption = select.options[select.selectedIndex];
                                        const selectedValue = select.value; // Obtiene la opción seleccionada
                                        if ((selectedValue) && selectedValue!="no") { // Solo incluir si existe una opción seleccionada
                                            filters[select.id] = selectedOption.textContent.trim(); // Agregar el texto de la opción
                                        }
                                    });

                                        // Obtener valores de los customSelects
                                    customSelects.forEach((customSelect) => {
                                        const selectElement = document.getElementById(customSelect.id);
                                        const optionsContainer = selectElement.querySelector('.options-container');
                                        
                                        // Verificar si hay checkboxes seleccionados
                                        const selectedOptions = Array.from(
                                            optionsContainer.querySelectorAll('.category-checkbox:checked')
                                        ).map(option => option.value);
                                        
                                        //Solo agregamos el filtro si hay opciones seleccionadas
                                        if (selectedOptions.length > 0) {
                                            filters[customSelect.id] = selectedOptions.join(', ');
                                        }
                                        // if (selectedOptions.length > 0) {
                                        //     const selectName = selectElement.name || customSelect.id; // Usa name si existe, si no, usa id
                                        //     filters[selectName] = selectedOptions.join(', '); // Agrega al filtro usando el name
                                        // }
                                    });



                                        // Actualizar el contenido de results
                                        const resultsDiv = document.getElementById('results');

                                        if (Object.keys(filters).length > 0) {
                                            // Convertir los filtros a texto legible
                                            const readableFilters = Object.entries(filters)
                                                .map(([key, value]) => `${key}: ${value}`) // Formato clave: valor
                                                .join('<br>'); // Separar cada filtro con una nueva línea

                                            resultsDiv.innerHTML = readableFilters; // Mostrar el texto formateado
                                        } else {
                                            resultsDiv.textContent = 'No hay resultados.'; // Mensaje si no hay filtros
                                        }
                                        

                              
                             });
  
        // fin activa funcionalidad de los filtros

    }
});
 

  