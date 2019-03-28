import {
  Address,
  arrayType,
  CommonResourceConf,
  Model,
  userPropsModel
} from '@liquid-labs/catalyst-core-api'

const personPropsModel = [
  'displayName',
  'phone',
  'email',
  'phoneBackup',
  'photoUrl']
  .map((propName) => ({ propName : propName, writable : true }))
personPropsModel.push(...userPropsModel)
personPropsModel.push({
  propName  : 'addresses',
  model     : Address,
  valueType : arrayType,
  writable  : true})
personPropsModel.push({
  propName            : 'changeDesc',
  unsetForNew         : true,
  writable            : true,
  optionalForComplete : true
})

const Person = class extends Model {
  get resourceName() { return 'persons' }
}
Model.finalizeConstructor(Person, personPropsModel)

const personResourceConf = new CommonResourceConf('person', {
  model       : Person,
  sortOptions : [
    { label : 'Dispaly name (asc)',
      value : 'displayName-asc',
      func  : (a, b) => a.displayName.localeCompare(b.displayName) },
    { label : 'Display name (desc)',
      value : 'displayName-desc',
      func  : (a, b) => -a.displayName.localeCompare(b.displayName) }
  ],
  sortDefault : 'displayName-asc'
})

export { Person, personPropsModel, personResourceConf }
