import { Model, Address, arrayType, entityPropModel } from '@liquid-labs/catalyst-core-api'

const personPropsModel = [
    'name',
    'phone',
    'email',
    'phoneBackup']
  .map((propName) => ({ propName: propName, writable: true}))
personPropsModel.push(...entityPropModel)
personPropsModel.push({
  propName: 'addresses',
  model: Address,
  valueType: arrayType,
  writable: true})
personPropsModel.push({
  propName: 'changeDesc',
  unsetForNew: true,
  writable: true,
  optionalForComplete: true
})

const Person = class extends Model {
  constructor(props = {}, opts = {}) {
    super(props, opts)
    if (Person === this.constructor) {
      throw new Error("Cannot create a Person directly, must first define a concrete sub-class.")
    }
  }
  // TODO: except sometimes it's a driver... where do we use this?
  get resourceName() { return 'persons' }
}
Model.finalizeConstructor(Person, personPropsModel)

const personResourceConf = new CommonResourceConf('person', {
  model: model.Person,
  sortOptions: [
    { label: 'Dispaly name (asc)',
      value: 'name-asc',
      func: (a, b) => a.displayName.localeCompare(b.displayName) },
    { label: 'Display name (desc)',
      value: 'name-desc',
      func: (a, b) => -a.displayName.localeCompare(b.displayName) }
  ],
  sortDefault: 'name-asc'
})

export { Person, personPropsModel, personResourceConf }
