import { Model, Address, arrayType, entityPropModel } from '@liquid-labs/catalyst-core-api'

export const personPropsModel = [
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

export default Person
